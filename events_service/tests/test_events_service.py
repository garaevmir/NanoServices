# import pytest
# from unittest.mock import AsyncMock, MagicMock
# from repository import PostRepository
# from events_service import PostService
# import datetime

# @pytest.fixture
# def mock_pool():
#     return MagicMock()

# @pytest.fixture
# def post_service(mock_pool):
#     return PostService(mock_pool)

# @pytest.mark.asyncio
# async def test_make_response():
#     mock_pool = MagicMock()
#     service = PostService(mock_pool)
    
#     test_post = {
#         "id": "post123",
#         "title": "Test",
#         "description": "Description",
#         "user_id": "user123",
#         "created_at": datetime.datetime.now(),
#         "updated_at": datetime.datetime.now(),
#         "is_private": False,
#         "tags": ["tag1"]
#     }
    
#     response = service.MakeResponse(test_post)
    
#     assert response["id"] == "post123"
#     assert response["title"] == "Test"
#     assert isinstance(response["created_at"], str)

# @pytest.mark.asyncio
# async def test_full_flow(post_service):
#     mock_repo = post_service.repo
#     mock_repo.create_post = AsyncMock(return_value={"id": "new_post"})
#     mock_repo.update_post = AsyncMock(return_value={"id": "updated_post"})
#     mock_repo.get_post = AsyncMock(return_value={"id": "found_post"})
#     mock_repo.delete_post = AsyncMock(return_value=True)

#     # Test create
#     created = await mock_repo.create_post(
#         title="Test", 
#         description="Desc", 
#         user_id="user1", 
#         is_private=False, 
#         tags=["tag1"]
#     )
#     assert created["id"] == "new_post"

#     # Test update
#     updated = await mock_repo.update_post(
#         "post123", 
#         "user1", 
#         {"title": "New Title"}
#     )
#     assert updated["id"] == "updated_post"

#     # Test get
#     found = await mock_repo.get_post("post123", "user1")
#     assert found["id"] == "found_post"

#     # Test delete
#     deleted = await mock_repo.delete_post("post123", "user1")
#     assert deleted is True