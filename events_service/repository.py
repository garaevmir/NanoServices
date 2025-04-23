import asyncpg
import sys

class PostRepository:
    def __init__(self, pool: asyncpg.Pool):
        self.pool = pool

    async def create_post(self, title: str, description: str, user_id: str, is_private: bool, tags: list):
        query = """
            INSERT INTO posts 
            (title, description, user_id, is_private, tags, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
            RETURNING *
        """
        return await self.pool.fetchrow(
            query, 
            title, description, user_id, is_private, tags
        )

    async def delete_post(self, post_id: str, user_id: str):
        query = "DELETE FROM posts WHERE id = $1 AND user_id = $2"
        result = await self.pool.execute(query, post_id, user_id)
        return "DELETE 1" in result

    async def update_post(self, post_id: str, user_id: str, fields):
        current_post = await self.get_post(post_id, user_id)
        if not current_post:
            return None

        merged_fields = {
            "title": current_post["title"],
            "description": current_post["description"],
            "is_private": current_post["is_private"],
            "tags": current_post["tags"]
        }
        merged_fields.update(fields)

        query = """
            UPDATE posts 
            SET title = $1, description = $2, is_private = $3, tags = $4, updated_at = NOW()
            WHERE id = $5 AND user_id = $6
            RETURNING *
        """
        
        return await self.pool.fetchrow(
            query,
            merged_fields["title"],
            merged_fields["description"],
            merged_fields["is_private"],
            merged_fields["tags"],
            post_id,
            user_id
        )

    async def get_post(self, post_id: str, user_id: str):
        query = """
            SELECT * FROM posts 
            WHERE id = $1 AND (NOT is_private OR user_id = $2)
        """
        return await self.pool.fetchrow(query, post_id, user_id)

    async def list_posts(self, page: int, page_size: int, user_id: str):
        offset = (page - 1) * page_size
        query = """
            SELECT * FROM posts 
            WHERE NOT is_private OR user_id = $3
            ORDER BY created_at DESC
            LIMIT $1 OFFSET $2
        """
        print("query", file=sys.stderr)
        posts = await self.pool.fetch(query, 
            page_size,
            offset,
            user_id
        )
        total = await self.pool.fetchval("SELECT COUNT(*) FROM posts")
        print("sended", file=sys.stderr)
        return posts, total
    
    async def add_comment(self, post_id: str, user_id: str, content: str):
        query = """
            INSERT INTO comments (post_id, user_id, content, created_at)
            VALUES ($1, $2, $3, NOW())
            RETURNING *
        """
        return await self.pool.fetchrow(query, post_id, user_id, content)

    async def get_comments(self, post_id: str, page: int, page_size: int):
        offset = (page - 1) * page_size
        query = """
            SELECT * FROM comments
            WHERE post_id = $1
            ORDER BY created_at DESC
            LIMIT $2 OFFSET $3
        """
        return await self.pool.fetch(query, post_id, page_size, offset)

    async def get_total_comments(self, post_id: str):
        return await self.pool.fetchval(
            "SELECT COUNT(*) FROM comments WHERE post_id = $1", post_id
        )

async def create_pool(dsn: str) -> asyncpg.Pool:
    return await asyncpg.create_pool(dsn=dsn)