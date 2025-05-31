import pytest
from unittest.mock import AsyncMock
from repository import PostgresManager
from datetime import datetime, timedelta
from generated import statistics_pb2

@pytest.mark.asyncio
async def test_get_post_stats():
    mock_pool = AsyncMock()
    mock_pool.fetchrow = AsyncMock(return_value=(100, 20, 5))
    db = PostgresManager(mock_pool)
    
    result = await db.get_post_stats("post123")
    
    assert result == (100, 20, 5)
    sql = mock_pool.fetchrow.call_args[0][0]
    assert "COUNT(*) FILTER (WHERE event_type = 'view')" in sql
    assert mock_pool.fetchrow.call_args[0][1] == "post123"

@pytest.mark.asyncio
async def test_get_trend():
    mock_pool = AsyncMock()
    db = PostgresManager(mock_pool)
    expected_days = 7
    start_date = (datetime.now() - timedelta(days=expected_days)).date()
    
    await db.get_trend("post123", "like", expected_days)
    
    sql = mock_pool.fetch.call_args[0][0]
    assert "GROUP BY DATE(event_time)" in sql
    assert mock_pool.fetch.call_args[0][1:] == ("post123", "like", start_date)

@pytest.mark.asyncio
async def test_get_top_posts():
    mock_pool = AsyncMock()
    db = PostgresManager(mock_pool)
    
    await db.get_top_posts(0)
    
    sql = mock_pool.fetch.call_args[0][0]
    assert "WHERE event_type = $1" in sql
    assert "ORDER BY count DESC" in sql
    assert mock_pool.fetch.call_args[0][1] == "view"

@pytest.mark.asyncio
async def test_get_top_users():
    mock_pool = AsyncMock()
    db = PostgresManager(mock_pool)
    
    await db.get_top_users(2)
    
    sql = mock_pool.fetch.call_args[0][0]
    assert "WHERE event_type = $1" in sql
    assert "GROUP BY user_id" in sql
    assert mock_pool.fetch.call_args[0][1] == "comment"

@pytest.mark.asyncio
async def test_get_trend_all_time():
    mock_pool = AsyncMock()
    db = PostgresManager(mock_pool)
    start_date = (datetime.now() - timedelta(days=365*10)).date()
    
    await db.get_trend("post123", "view", 365*10)
    
    assert mock_pool.fetch.call_args[0][3] == start_date