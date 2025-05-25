import pytest
from unittest.mock import AsyncMock, MagicMock
from grpc import StatusCode
import datetime
from stats_server import StatsService
from repository import PostgresManager
from generated import statistics_pb2

@pytest.fixture
def mock_db():
    return MagicMock(spec=PostgresManager)

@pytest.fixture
def stats_service(mock_db):
    return StatsService(mock_db)

@pytest.fixture
def mock_context():
    context = MagicMock()
    context.abort = AsyncMock()
    context.set_code = MagicMock()
    context.set_details = MagicMock()
    return context

@pytest.mark.asyncio
async def test_get_post_stats_success(stats_service, mock_db, mock_context):
    mock_db.get_post_stats = AsyncMock(return_value=(100, 20, 5))
    request = statistics_pb2.PostStatsRequest(post_id="post123")
    
    response = await stats_service.GetPostStats(request, mock_context)
    
    assert response.views == 100
    assert response.likes == 20
    assert response.comments == 5
    mock_db.get_post_stats.assert_called_once_with("post123")

@pytest.mark.asyncio
async def test_get_post_stats_error(stats_service, mock_db, mock_context):
    mock_db.get_post_stats = AsyncMock(side_effect=Exception("DB error"))
    
    await stats_service.GetPostStats(statistics_pb2.PostStatsRequest(), mock_context)
    
    mock_context.abort.assert_called_with(StatusCode.INTERNAL, "DB error")

@pytest.mark.asyncio
async def test_get_trend_success(stats_service, mock_db):
    test_data = [
        {'date': datetime.date(2023, 1, 1), 'count': 10},
        {'date': datetime.date(2023, 1, 2), 'count': 15}
    ]
    mock_db.get_trend = AsyncMock(return_value=test_data)
    request = statistics_pb2.PostTrendRequest(post_id="post123", period="7d")
    
    response = await stats_service.GetViewsTrend(request, None)
    
    assert len(response.data) == 2
    assert response.data[0].date == "2023-01-01"
    assert response.data[1].count == 15
    mock_db.get_trend.assert_called_with("post123", "view", 7)

@pytest.mark.asyncio
async def test_get_trend_all_period(stats_service, mock_db):
    mock_db.get_trend = AsyncMock()
    request = statistics_pb2.PostTrendRequest(post_id="post123", period="all")
    
    await stats_service.GetLikesTrend(request, None)
    
    mock_db.get_trend.assert_called_with("post123", "like", 3650)

@pytest.mark.asyncio
async def test_get_trend_invalid_period(stats_service, mock_db, mock_context):
    request = statistics_pb2.PostTrendRequest(post_id="post123", period="invalid")
    
    await stats_service.GetCommentsTrend(request, mock_context)
    
    mock_context.abort.assert_called_with(StatusCode.INTERNAL, "invalid literal for int() with base 10: 'invali'")

@pytest.mark.asyncio
async def test_get_top_posts_success(stats_service, mock_db):
    test_data = [{"post_id": "p1", "count": 100}, {"post_id": "p2", "count": 80}]
    mock_db.get_top_posts = AsyncMock(return_value=test_data)
    request = statistics_pb2.TopRequest(metric=statistics_pb2.TopRequest.Metric.VIEWS)
    
    response = await stats_service.GetTopPosts(request, None)
    
    assert len(response.posts) == 2
    assert response.posts[0].post_id == "p1"
    assert response.posts[1].count == 80

@pytest.mark.asyncio
async def test_get_top_posts_invalid_metric(stats_service, mock_db, mock_context):
    mock_db.get_top_posts = AsyncMock(side_effect=ValueError("Invalid metric: 5"))
    request = statistics_pb2.TopRequest(metric=5)
    
    await stats_service.GetTopPosts(request, mock_context)
    
    mock_context.set_code.assert_called_with(StatusCode.INVALID_ARGUMENT)
    mock_context.set_details.assert_called_with("Invalid metric: 5")

@pytest.mark.asyncio
async def test_get_top_users_success(stats_service, mock_db):
    test_data = [{"user_id": "u1", "count": 50}, {"user_id": "u2", "count": 30}]
    mock_db.get_top_users = AsyncMock(return_value=test_data)
    request = statistics_pb2.TopRequest(metric=statistics_pb2.TopRequest.Metric.LIKES)
    
    response = await stats_service.GetTopUsers(request, None)
    
    assert len(response.users) == 2
    assert response.users[0].user_id == "u1"
    assert response.users[1].count == 30

@pytest.mark.asyncio
async def test_get_top_users_error(stats_service, mock_db, mock_context):
    mock_db.get_top_users = AsyncMock(side_effect=Exception("DB error"))
    request = statistics_pb2.TopRequest(metric=statistics_pb2.TopRequest.Metric.COMMENTS)
    
    await stats_service.GetTopUsers(request, mock_context)
    
    mock_context.set_code.assert_called_with(StatusCode.INTERNAL)
    mock_context.set_details.assert_called_with("Internal error: DB error")