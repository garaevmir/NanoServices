import datetime
import pytest
from unittest.mock import AsyncMock, MagicMock
from grpc import StatusCode
import asyncpg

from events_server import PostService
from generated import post_pb2

@pytest.fixture
def mock_pool():
    return MagicMock(spec=asyncpg.Pool)

@pytest.fixture
def mock_kafka_producer():
    return MagicMock()

@pytest.fixture
def post_service(mock_pool, mock_kafka_producer):
    return PostService(mock_pool, mock_kafka_producer)

@pytest.fixture
def mock_context():
    context = MagicMock()
    context.set_code = MagicMock()
    context.set_details = MagicMock()
    return context

@pytest.mark.asyncio
async def test_create_post_success(post_service, mock_context):
    mock_post = {
        "id": "1",
        "title": "Test",
        "description": "Desc",
        "user_id": "user1",
        "created_at": datetime.datetime(2023, 10, 1, 0, 0, 0),
        "updated_at": datetime.datetime(2023, 10, 1, 0, 0, 0),
        "is_private": False,
        "tags": ["tag1"]
    }
    post_service.repo.create_post = AsyncMock(return_value=mock_post)
    request = post_pb2.CreatePostRequest(
        title="Test",
        description="Desc",
        user_id="user1",
        is_private=False,
        tags=["tag1"]
    )
    response = await post_service.CreatePost(request, mock_context)
    assert response.id == "1"
    assert response.title == "Test"

@pytest.mark.asyncio
async def test_update_post_success(post_service, mock_context):
    mock_post = {
        "id": "1",
        "title": "New Title",
        "description": "New Desc",
        "created_at": datetime.datetime(2023, 10, 1, 0, 0, 0),
        "updated_at": datetime.datetime(2023, 10, 1, 0, 0, 0),
        "is_private": True,
        "tags": ["new_tag"],
        "user_id": "user1"
    }
    post_service.repo.update_post = AsyncMock(return_value=mock_post)
    request = post_pb2.UpdatePostRequest(
        post_id="1",
        user_id="user1",
        title="New Title",
        is_private=True
    )
    response = await post_service.UpdatePost(request, mock_context)
    assert response.title == "New Title"
    assert response.is_private is True

@pytest.mark.asyncio
async def test_get_post_success(post_service, mock_context):
    mock_post = {
        "id": "1",
        "title": "Test",
        "description": "Desc",
        "user_id": "user1",
        "created_at": datetime.datetime(2023, 10, 1, 0, 0, 0),
        "updated_at": datetime.datetime(2023, 10, 1, 0, 0, 0),
        "is_private": False,
        "tags": []
    }
    post_service.repo.get_post = AsyncMock(return_value=mock_post)
    request = post_pb2.GetPostRequest(post_id="1", user_id="user1")
    response = await post_service.GetPost(request, mock_context)
    assert response.id == "1"

@pytest.mark.asyncio
async def test_list_posts_success(post_service, mock_context):
    mock_posts = [
        {
            "id": "1",
            "title": "Post1",
            "description": "Desc1",
            "user_id": "user1",
            "created_at": datetime.datetime(2023, 10, 1, 0, 0, 0),
            "updated_at": datetime.datetime(2023, 10, 1, 0, 0, 0),
            "is_private": False,
            "tags": []
        },
        {
            "id": "2",
            "title": "Post2",
            "description": "Desc2",
            "user_id": "user1",
            "created_at": datetime.datetime(2023, 10, 1, 0, 0, 0),
            "updated_at": datetime.datetime(2023, 10, 1, 0, 0, 0),
            "is_private": False,
            "tags": []
        }
    ]
    post_service.repo.list_posts = AsyncMock(return_value=(mock_posts, 2))
    request = post_pb2.ListPostsRequest(page=1, page_size=10, user_id="user1")
    response = await post_service.ListPosts(request, mock_context)
    assert len(response.posts) == 2
    assert response.total == 2

@pytest.mark.asyncio
async def test_create_post_exception(post_service, mock_context):
    post_service.repo.create_post = AsyncMock(side_effect=Exception("DB error"))
    request = post_pb2.CreatePostRequest()
    
    response = await post_service.CreatePost(request, mock_context)
    
    mock_context.set_code.assert_called_with(StatusCode.INTERNAL)
    assert response == post_pb2.PostResponse()

@pytest.mark.asyncio
async def test_delete_post_failure(post_service, mock_context):
    post_service.repo.delete_post = AsyncMock(return_value=False)
    request = post_pb2.DeletePostRequest(post_id="1", user_id="user1")
    
    response = await post_service.DeletePost(request, mock_context)
    
    mock_context.set_code.assert_called_with(StatusCode.NOT_FOUND)
    assert response.success is False

@pytest.mark.asyncio
async def test_update_post_full_fields(post_service, mock_context):
    mock_post = {
        "id": "1",
        "title": "New",
        "description": "New Desc",
        "is_private": True,
        "tags": ["new_tag"],
        "user_id": "user1",
        "created_at": datetime.datetime.now(),
        "updated_at": datetime.datetime.now()
    }
    post_service.repo.update_post = AsyncMock(return_value=mock_post)
    
    request = post_pb2.UpdatePostRequest(
        post_id="1",
        user_id="user1",
        title="New",
        description="New Desc",
        is_private=True,
        tags=["new_tag"]
    )
    
    response = await post_service.UpdatePost(request, mock_context)
    assert response.description == "New Desc"
    assert response.tags == ["new_tag"]

@pytest.mark.asyncio
async def test_update_post_not_found(post_service, mock_context):
    post_service.repo.update_post = AsyncMock(return_value=None)
    request = post_pb2.UpdatePostRequest(post_id="999", user_id="user1")
    
    response = await post_service.UpdatePost(request, mock_context)
    
    mock_context.set_code.assert_called_with(StatusCode.NOT_FOUND)
    assert response == post_pb2.PostResponse()

@pytest.mark.asyncio
async def test_get_post_not_found(post_service, mock_context):
    post_service.repo.get_post = AsyncMock(return_value=None)
    request = post_pb2.GetPostRequest(post_id="999", user_id="user1")
    
    response = await post_service.GetPost(request, mock_context)
    
    mock_context.set_code.assert_called_with(StatusCode.NOT_FOUND)
    assert response == post_pb2.PostResponse()

@pytest.mark.asyncio
async def test_list_posts_pagination(post_service, mock_context):
    post_service.repo.list_posts = AsyncMock(return_value=([], 0))
    
    request = post_pb2.ListPostsRequest(page=0, page_size=5, user_id="user1")
    await post_service.ListPosts(request, mock_context)
    post_service.repo.list_posts.assert_called_with(page=1, page_size=5, user_id="user1")
    
    request = post_pb2.ListPostsRequest(page=2, page_size=150, user_id="user1")
    await post_service.ListPosts(request, mock_context)
    post_service.repo.list_posts.assert_called_with(page=2, page_size=10, user_id="user1")

@pytest.mark.asyncio
async def test_list_posts_exception(post_service, mock_context):
    post_service.repo.list_posts = AsyncMock(side_effect=Exception("DB error"))
    request = post_pb2.ListPostsRequest()
    
    response = await post_service.ListPosts(request, mock_context)
    
    mock_context.set_code.assert_called_with(StatusCode.INTERNAL)
    assert response == post_pb2.ListPostsResponse()

@pytest.mark.asyncio
async def test_update_post_exception(post_service, mock_context):
    post_service.repo.update_post = AsyncMock(side_effect=Exception("DB Error"))
    
    request = post_pb2.UpdatePostRequest(
        post_id="1",
        user_id="user1",
        title="New Title"
    )
    
    response = await post_service.UpdatePost(request, mock_context)
    
    mock_context.set_code.assert_called_with(StatusCode.INTERNAL)
    mock_context.set_details.assert_called_with("Error: DB Error")
    assert response == post_pb2.PostResponse()

@pytest.mark.asyncio
async def test_view_post_success(post_service, mock_context):
    post_service._send_kafka_event = AsyncMock()
    request = post_pb2.ViewPostRequest(post_id="1", user_id="user1")
    
    response = await post_service.ViewPost(request, mock_context)
    
    assert response.success is True
    post_service._send_kafka_event.assert_called_with(
        "post_views", "user1", "1"
    )

@pytest.mark.asyncio
async def test_like_post_success(post_service, mock_context):
    post_service._send_kafka_event = AsyncMock()
    request = post_pb2.LikePostRequest(post_id="1", user_id="user1")
    
    response = await post_service.LikePost(request, mock_context)
    
    assert response.success is True
    post_service._send_kafka_event.assert_called_with(
        "post_likes", "user1", "1"
    )

@pytest.mark.asyncio
async def test_comment_post_success(post_service, mock_context):
    mock_comment = {"id": "c1", "content": "test"}
    post_service.repo.add_comment = AsyncMock(return_value=mock_comment)
    post_service._send_kafka_event = AsyncMock()
    
    request = post_pb2.CommentPostRequest(
        post_id="1", 
        user_id="user1", 
        content="test"
    )
    
    response = await post_service.CommentPost(request, mock_context)
    
    assert response.comment_id == "c1"
    post_service._send_kafka_event.assert_called_with(
        "post_comments", "user1", "1", "test"
    )

@pytest.mark.asyncio
async def test_get_comments_success(post_service, mock_context):
    mock_comments = [
        {"id": "c1", "content": "test1", "user_id": "user1", "created_at": datetime.datetime.now()},
        {"id": "c2", "content": "test2", "user_id": "user2", "created_at": datetime.datetime.now()}
    ]
    post_service.repo.get_comments = AsyncMock(return_value=mock_comments)
    post_service.repo.get_total_comments = AsyncMock(return_value=2)
    
    request = post_pb2.GetCommentsRequest(
        post_id="1", 
        page=1, 
        page_size=10, 
        user_id="user1"
    )
    
    response = await post_service.GetComments(request, mock_context)
    
    assert len(response.comments) == 2
    assert response.total == 2

@pytest.mark.asyncio
async def test_send_kafka_event_error(post_service, mock_context):
    post_service._send_kafka_event = AsyncMock(side_effect=Exception("Kafka error"))
    request = post_pb2.ViewPostRequest(post_id="1", user_id="user1")
    
    response = await post_service.ViewPost(request, mock_context)
    
    assert response.success is True
    post_service._send_kafka_event.assert_called_with(
        "post_views", "user1", "1"
    )