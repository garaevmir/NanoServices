from datetime import datetime, timedelta
import sys
import asyncpg

class PostgresManager:
    def __init__(self, pool: asyncpg.Pool):
        self.pool = pool
        
    async def insert_event(self, event: dict):
        query = """
            INSERT INTO events (
                event_time, 
                event_type, 
                post_id, 
                user_id, 
                content
            ) VALUES ($1, $2, $3, $4, $5)
        """
        
        await self.pool.execute(
            query,
            event['event_time'],
            event['event_type'],
            event['post_id'],
            event['user_id'],
            event.get('content')
        )

    async def get_post_stats(self, post_id: str) -> tuple:
        query = """
            SELECT 
                COUNT(*) FILTER (WHERE event_type = 'view') as views,
                COUNT(*) FILTER (WHERE event_type = 'like') as likes,
                COUNT(*) FILTER (WHERE event_type = 'comment') as comments
            FROM events
            WHERE post_id = $1
        """
        return await self.pool.fetchrow(
            query, 
            post_id,
        )


    async def get_trend(self, post_id: str, metric: str, days: int) -> list:
        start_date = (datetime.now() - timedelta(days=days)).date()
        print(start_date, file=sys.stderr)
        query = """
            SELECT 
                DATE(event_time) as date,
                COUNT(*) as count
            FROM events
            WHERE 
                post_id = $1 AND
                event_type = $2 AND
                event_time >= $3
            GROUP BY DATE(event_time)
            ORDER BY DATE(event_time)
        """
        result = await self.pool.fetch(
            query, 
            post_id,
            metric,
            start_date
        )
        
        return [{
            "date": row[0], 
            "count": row[1]
        } for row in result]
        
    async def get_top_posts(self, metric: str) -> list:
        metric_map = {
            0: "view",
            1: "like",
            2: "comment"
        }
        event_type = metric_map.get(metric)
        
        if not event_type:
            raise ValueError(f"Invalid metric: {metric}")
        
        query = """
            SELECT 
                post_id,
                COUNT(*) as count
            FROM events
            WHERE event_type = $1
            GROUP BY post_id
            ORDER BY count DESC
            LIMIT 10
        """
        result = await self.pool.fetch(
            query, 
            event_type
        )
        
        return [{
            "post_id": row["post_id"],
            "count": row["count"]
        } for row in result]

    async def get_top_users(self, metric: str) -> list:
        metric_map = {
            0: "view",
            1: "like",
            2: "comment"
        }
        event_type = metric_map.get(metric)
        
        if not event_type:
            raise ValueError(f"Invalid metric: {metric}")
        
        query = """
            SELECT 
                user_id,
                COUNT(*) as count
            FROM events
            WHERE event_type = $1
            GROUP BY user_id
            ORDER BY count DESC
            LIMIT 10
        """
        result = await self.pool.fetch(
            query, 
            event_type
        )
        
        return [{
            "user_id": row["user_id"],
            "count": row["count"]
        } for row in result]
        
        
async def create_pool(dsn: str) -> asyncpg.Pool:
    return await asyncpg.create_pool(dsn=dsn)