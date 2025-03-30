import grpc
from grpc import aio
from generated import post_pb2_grpc, post_pb2
from repository import PostRepository, create_pool

class PostService(post_pb2_grpc.PostServiceServicer):
    def __init__(self, pool):
        self.repo = PostRepository(pool)

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
    pool = await create_pool(
        "postgres://postgres:postgres@events_db:5432/events_db"
    )
    
    server = aio.server()
    post_pb2_grpc.add_PostServiceServicer_to_server(
        PostService(pool), server)
    server.add_insecure_port('[::]:50051')
    await server.start()
    await server.wait_for_termination()

if __name__ == '__main__':
    import asyncio
    asyncio.run(serve())