import pytest
from unittest.mock import AsyncMock
from repository import PostRepository

@pytest.mark.asyncio
async def test_create_post():
    mock_pool = AsyncMock()
    repo = PostRepository(mock_pool)
    
    test_data = {
        "title": "Test Post",
        "description": "Test Description",
        "user_id": "user123",
        "is_private": False,
        "tags": ["tag1", "tag2"]
    }
    
    await repo.create_post(**test_data)
    
    sql = mock_pool.fetchrow.call_args[0][0]
    assert "INSERT INTO posts" in sql
    assert mock_pool.fetchrow.call_args[0][1:] == (
        test_data["title"],
        test_data["description"],
        test_data["user_id"],
        test_data["is_private"],
        test_data["tags"]
    )

@pytest.mark.asyncio
async def test_get_post():
    mock_pool = AsyncMock()
    repo = PostRepository(mock_pool)
    
    await repo.get_post("post123", "user123")
    
    sql = mock_pool.fetchrow.call_args[0][0]
    assert "SELECT * FROM posts" in sql
    assert "NOT is_private OR user_id = $2" in sql
    assert mock_pool.fetchrow.call_args[0][1:] == ("post123", "user123")

@pytest.mark.asyncio
async def test_update_post():
    mock_pool = AsyncMock()
    repo = PostRepository(mock_pool)
    
    update_fields = {
        "title": "New Title",
        "description": "New Description",
        "is_private": True,
        "tags": ["new_tag"]
    }
    
    await repo.update_post("post123", "user123", update_fields)
    
    sql = mock_pool.fetchrow.call_args[0][0]
    assert "UPDATE posts" in sql
    assert mock_pool.fetchrow.call_args[0][1:-2] == (
        update_fields["title"],
        update_fields["description"],
        update_fields["is_private"],
        update_fields["tags"]
    )

@pytest.mark.asyncio
async def test_list_posts():
    mock_pool = AsyncMock()
    mock_pool.fetch.return_value = [{"id": "1"}, {"id": "2"}]
    mock_pool.fetchval.return_value = 10
    repo = PostRepository(mock_pool)
    
    posts, total = await repo.list_posts(page=2, page_size=5, user_id="user123")
    
    assert len(posts) == 2
    assert total == 10
    assert "LIMIT $1 OFFSET $2" in mock_pool.fetch.call_args[0][0]
    assert mock_pool.fetch.call_args[0][1:3] == (5, 5)

@pytest.mark.asyncio
async def test_delete_post():
    mock_pool = AsyncMock()
    mock_pool.execute.return_value = "DELETE 1"
    repo = PostRepository(mock_pool)
    
    result = await repo.delete_post("post123", "user123")
    
    sql = mock_pool.execute.call_args[0][0]
    assert "DELETE FROM posts WHERE id = $1 AND user_id = $2" in sql
    assert mock_pool.execute.call_args[0][1:] == ("post123", "user123")
    
    assert result is True

@pytest.mark.asyncio
async def test_delete_post_not_found():
    mock_pool = AsyncMock()
    mock_pool.execute.return_value = "DELETE 0"
    repo = PostRepository(mock_pool)
    
    result = await repo.delete_post("invalid_id", "user123")
    
    assert result is False