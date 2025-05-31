import asyncio
import os
import json
import logging
from datetime import datetime
import sys

from aiokafka import AIOKafkaConsumer
import grpc
from grpc import aio
from repository import PostgresManager, create_pool
from generated import statistics_pb2_grpc, statistics_pb2

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("StatsService")

class StatsService(statistics_pb2_grpc.StatsServiceServicer):
    def __init__(self, db):
        self.db = db

    async def GetPostStats(self, request, context):
        try:
            views, likes, comments = await self.db.get_post_stats(request.post_id)
            return statistics_pb2.PostStatsResponse(
                views=views or 0,
                likes=likes or 0,
                comments=comments or 0
            )
        except Exception as e:
            logger.error(f"GetPostStats error: {str(e)}")
            await context.abort(grpc.StatusCode.INTERNAL, str(e))

    async def GetViewsTrend(self, request, context):
        return await self._handle_trend(request, context, 'view')

    async def GetLikesTrend(self, request, context):
        return await self._handle_trend(request, context, 'like')

    async def GetCommentsTrend(self, request, context):
        return await self._handle_trend(request, context, 'comment')

    async def _handle_trend(self, request, context, metric):
        try:
            days = int(request.period[:-1]) if request.period != 'all' else 365*10
            data = await self.db.get_trend(request.post_id, metric, days)
            return statistics_pb2.PostTrendResponse(
                data=[statistics_pb2.TrendItem(
                    date=d['date'].isoformat(),
                    count=d['count']
                ) for d in data]
            )
        except Exception as e:
            logger.error(f"GetTrend error: {str(e)}")
            await context.abort(grpc.StatusCode.INTERNAL, str(e))
                
    async def GetTopPosts(self, request, context):
        try:
            print('start getting top posts', file=sys.stderr)
            
            data = await self.db.get_top_posts(request.metric)
            return statistics_pb2.TopPostsResponse(
                posts=[statistics_pb2.PostItem(
                    post_id=item["post_id"],
                    count=item["count"]
                ) for item in data]
            )
        except ValueError as e:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            print(str(e), file=sys.stderr)
            context.set_details(str(e))
            return statistics_pb2.TopPostsResponse()
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            print(str(e), file=sys.stderr)
            context.set_details(str(e))
            return statistics_pb2.TopPostsResponse()

    async def GetTopUsers(self, request, context):
        try:
            
            data = await self.db.get_top_users(request.metric)
            return statistics_pb2.TopUsersResponse(
                users=[statistics_pb2.UserItem(
                    user_id=item["user_id"],
                    count=item["count"]
                ) for item in data]
            )
        except ValueError as e:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(str(e))
            return statistics_pb2.TopUsersResponse()
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Internal error: {str(e)}")
            return statistics_pb2.TopUsersResponse()
        
async def consume_events(db):
    consumer = AIOKafkaConsumer(
        'post_views',
        'post_likes',
        'post_comments',
        bootstrap_servers='kafka:9092',
        group_id="stats-service-group",
        auto_offset_reset='earliest',
        max_poll_interval_ms=300000,
        session_timeout_ms=10000,
        request_timeout_ms=15000
    )
    
    event_type_mapping = {
        'post_views': 'view',
        'post_likes': 'like',
        'post_comments': 'comment'
    }
    try:
        await consumer.start()
        logger.info("Connected to Kafka")
        print("kafka started", file=sys.stderr)
        
        async for msg in consumer:
            try:
                event = json.loads(msg.value.decode())
                event_type = event_type_mapping[msg.topic]
                
                await db.insert_event({
                    'event_time': datetime.fromisoformat(event['timestamp']),
                    'event_type': event_type,
                    'post_id': event['post_id'],
                    'user_id': event['user_id'],
                    'content': event.get('content', '')
                })
                logger.debug(f"Processed event: {event['post_id']}")
            except Exception as e:
                logger.error(f"Error processing event: {str(e)}")
    finally:
        await consumer.stop()
        logger.info("Kafka consumer stopped")
        os._exit(0)
   

async def serve():
    pool = await create_pool(os.getenv('DATABASE_URL'))
    db = PostgresManager(pool)
    server = aio.server()
    statistics_pb2_grpc.add_StatsServiceServicer_to_server(StatsService(db), server)
    server.add_insecure_port('[::]:50052')
    await server.start()
    
    consumer_task = asyncio.create_task(consume_events(db))
    
    try:
        await server.wait_for_termination()
    except KeyboardInterrupt:
        await server.stop(5)

if __name__ == '__main__':
    asyncio.run(serve())