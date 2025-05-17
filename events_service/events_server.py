from datetime import datetime
import json
import sys
from aiokafka import AIOKafkaProducer
import grpc
from grpc import aio
from generated import post_pb2_grpc, post_pb2
from repository import PostRepository, create_pool
import json
import logging

logger = logging.getLogger(__name__)

class PostService(post_pb2_grpc.PostServiceServicer):
    def __init__(self, pool, kafka_producer):
        self.repo = PostRepository(pool)
        self.kafka = kafka_producer

    async def ViewPost(self, request, context):
        logger.info("ViewPost request for post_id: %s", request.post_id)
        try:
            await self._send_kafka_event("post_views", request.user_id, request.post_id)
        except Exception as e:
            logger.error("Kafka error: %s", str(e))
        return post_pb2.InteractionResponse(success=True)

    async def LikePost(self, request, context):
        logger.info("LikePost request for post_id: %s", request.post_id)
        try:
            await self._send_kafka_event("post_likes", request.user_id, request.post_id)
        except Exception as e:
            logger.error("Kafka error: %s", str(e))
        return post_pb2.InteractionResponse(success=True)

    async def CommentPost(self, request, context):
        try:
            comment = await self.repo.add_comment(
                post_id=request.post_id,
                user_id=request.user_id,
                content=request.content
            )
            await self._send_kafka_event(
                "post_comments", 
                request.user_id, 
                request.post_id,
                request.content
            )
            return post_pb2.CommentResponse(comment_id=str(comment['id']))
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return post_pb2.CommentResponse()

    async def GetComments(self, request, context):
        try:
            comments = await self.repo.get_comments(
                post_id=request.post_id,
                page=request.page,
                page_size=request.page_size
            )
            total = await self.repo.get_total_comments(request.post_id)
            return post_pb2.CommentsResponse(
                comments=[self._format_comment(c) for c in comments],
                total=total
            )
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return post_pb2.CommentsResponse()

    def _format_comment(self, comment):
        return post_pb2.Comment(
            id=str(comment['id']),
            content=comment['content'],
            user_id=str(comment['user_id']),
            created_at=comment['created_at'].isoformat()
        )

    async def _send_kafka_event(self, topic, user_id, post_id, content=None):
        logger.info('sending event to ', topic)
        event = {
            "user_id": user_id,
            "post_id": post_id,
            "timestamp": datetime.now().isoformat(),
            "content": content
        }
        logger.info(
            "Preparing to send event to Kafka. Topic: %s, Event: %s",
            topic,
            {k: v for k, v in event.items() if k != "content"}
        )
        await self.kafka.send(topic, json.dumps(event).encode())

    async def CreatePost(self, request, context):
        try:
            post = await self.repo.create_post(
                title=request.title,
                description=request.description,
                user_id=request.user_id,
                is_private=request.is_private,
                tags=list(request.tags)
            )
            
            return self.MakeResponse(post)
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return post_pb2.PostResponse()

    async def DeletePost(self, request, context):
        success = await self.repo.delete_post(request.post_id, request.user_id)
        if not success:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("Post not found or permission denied")
        return post_pb2.DeletePostResponse(success=success)

    async def UpdatePost(self, request, context):
        update_fields = {}
        
        if request.HasField("title"): 
            update_fields["title"] = request.title
        if request.HasField("description"): 
            update_fields["description"] = request.description
        if request.HasField("is_private"): 
            update_fields["is_private"] = request.is_private
        if request.tags:
            update_fields["tags"] = list(request.tags)

        try:
            post = await self.repo.update_post(
                post_id=request.post_id,
                user_id=request.user_id,
                fields=update_fields
            )
            
            if not post:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                return post_pb2.PostResponse()
                
            return self.MakeResponse(post)
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Error: {str(e)}")
            return post_pb2.PostResponse()

    async def GetPost(self, request, context):
        post = await self.repo.get_post(request.post_id, request.user_id)
        if not post:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            return post_pb2.PostResponse()
            
        return self.MakeResponse(post)

    async def ListPosts(self, request, context):
        try:
            page = int(request.page) if request.page > 0 else 1
            page_size = int(request.page_size) if 1 <= request.page_size <= 100 else 10
            
            posts, total = await self.repo.list_posts(
                page=page,
                page_size=page_size,
                user_id=request.user_id
            )
            
            return post_pb2.ListPostsResponse(
                posts=[self.MakeResponse(p) for p in posts],
                total=total
            )
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Error: {str(e)}")
            return post_pb2.ListPostsResponse()

    def MakeResponse(self, post):
        return post_pb2.PostResponse(
            id=str(post['id']),
            title=post['title'],
            description=post['description'],
            user_id=str(post['user_id']),
            created_at=post['created_at'].isoformat(),
            updated_at=post['updated_at'].isoformat(),
            is_private=post['is_private'],
            tags=post['tags'])


async def serve():
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
        stream=sys.stdout
    )
    logger.info("Starting gRPC server...")

    try:
        logger.info("Connecting to PostgreSQL...")
        pool = await create_pool("postgres://postgres:postgres@events_db:5432/events_db")
    except Exception as e:
        logger.error("PostgreSQL connection failed: %s", str(e))
        return

    producer = AIOKafkaProducer(
        bootstrap_servers='kafka:9092',
        request_timeout_ms=30000
    )
    try:
        logger.info("Connecting to Kafka...")
        await producer.start()
    except Exception as e:
        logger.error("Kafka connection failed: %s", str(e))
        return

    server = aio.server()
    post_pb2_grpc.add_PostServiceServicer_to_server(
        PostService(pool, producer), server
    )
    server.add_insecure_port('[::]:50051')
    logger.info("gRPC server started on port 50051")
    await server.start()
    await server.wait_for_termination()

if __name__ == '__main__':
    import asyncio
    asyncio.run(serve())