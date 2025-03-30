import pytest
from unittest.mock import patch

@pytest.fixture(autouse=True)
def mock_db():
    with patch("asyncpg.create_pool") as mock_pool:
        yield mock_pool