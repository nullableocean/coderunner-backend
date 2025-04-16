import base64
import pytest
import requests
import uuid
import time
import os

BASE_URL = os.getenv('BASE_URL', 'http://127.0.0.1:8080')
REGISTER_URL = BASE_URL + "/user/register"
LOGIN_URL = BASE_URL + "/user/login"
TASK_URL = BASE_URL + "/task"
STATUS_URL = BASE_URL + "/status"
RESULT_URL = BASE_URL + "/result"

@pytest.fixture(scope='module')
def user_data():
    login = f'user_{uuid.uuid4()}'
    password = 'password1234'
    return {'login': login, 'password': password}

@pytest.fixture(scope='module')
def auth_token(user_data):
    response = requests.post(LOGIN_URL, json=user_data)
    print("response auth token: ", response.text)
    assert response.status_code == 200

    data = response.json()
    assert 'token' in data

    return data['token']

def test_register_user(user_data):
    response = requests.post(REGISTER_URL, json=user_data)
    assert response.status_code == 201

def test_login_user(user_data):
    response = requests.post(LOGIN_URL, json=user_data)

    print("response login: ", response.text)
    assert response.status_code == 200
    data = response.json()
    assert 'token' in data

def get_code_processor_payload():
    return {"compiler": "python3", "code": "print('Hello, stdout world!')"}

def test_create_task(auth_token):
    headers = {'Authorization': f'Bearer {auth_token}'}

    payload = get_code_processor_payload()
    response = requests.post(TASK_URL, headers=headers, json=payload) 

    print("response create task: ", response.text)
    assert response.status_code == 201
    data = response.json()
    assert 'task_id' in data

    return data['task_id']

def test_task_status_and_result(auth_token):
    task_id = test_create_task(auth_token)
    status_url = f"{STATUS_URL}/{task_id}"
    result_url = f"{RESULT_URL}/{task_id}"
    headers = {'Authorization': f'Bearer {auth_token}'}

    retry = 10
    while retry >= 0:
        response = requests.get(status_url, headers=headers)
        assert response.status_code == 200
        data = response.json()
        assert 'status' in data
        
        if data['status'] == 'ready':
            break
         
        assert data['status'] == 'in_progress', f'undefined status: {data['status']}!'
        retry -= 1
        time.sleep(3)
    assert retry > 0, "task is still in progress!"

    response = requests.get(result_url, headers=headers)

    print("response result: ", response.text)

    assert response.status_code == 200
    data = response.json()
    assert 'result' in data

def test_task_not_found(auth_token):
    invalid_task_id = str(uuid.uuid4())
    status_url = f"{STATUS_URL}/{invalid_task_id}"
    result_url = f"{RESULT_URL}/{invalid_task_id}"
    headers = {'Authorization': f'Bearer {auth_token}'}

    response = requests.get(status_url, headers=headers)

    print("response task not found: ", response.text)
    assert response.status_code == 404

    response = requests.get(result_url, headers=headers)

    print("response result not found: ", response.text)
    assert response.status_code == 404

def test_unauthorized_access():
    invalid_task_id = str(uuid.uuid4())
    status_url = f"{STATUS_URL}/{invalid_task_id}"
    result_url = f"{RESULT_URL}/{invalid_task_id}"

    response = requests.post(TASK_URL)
    assert response.status_code == 401

    response = requests.get(status_url)
    assert response.status_code == 401

    response = requests.get(result_url)
    assert response.status_code == 401