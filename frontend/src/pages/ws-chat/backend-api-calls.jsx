import axios from "axios";

export const API_BACK_URL = "http://localhost:8080"

export const myAxios = axios.create({
    baseURL: API_BACK_URL,
    headers: {
        "Content-Type": "application/json",
    },
    withCredentials: true,
});

export const myAxiosWithoutHeaders = axios.create({
    baseURL: API_BACK_URL,
    withCredentials: true,
})

const sendUsersOnlineRequest = () => {
    const response = myAxiosWithoutHeaders.get("usersonline");
    return response
};

const sendLoginRequest = async (login, password) => {
    return await myAxios.post(
        "signin",
        JSON.stringify({
            username: login,
            password: password,
        })
    );
};

const sendLogoutRequest = async () => {
    return await myAxios.post(
        "signout",
        JSON.stringify({})
    );
}

const sendProfileRequest = () => {
    const response = myAxios.get("/me")
    return response
};

const sendChatMessagesRequest = (userToChat) => {
    const response = myAxiosWithoutHeaders.get(`/chat/${userToChat}?offset=0`)
    return response
}

const sendGroupChatMessagesRequest = (groupId) => {
    const response = myAxiosWithoutHeaders.get(`/group-chat/${groupId}?offset=0`)
    return response
}

const sendChatMessagesRequestWithOffset = (userToChat, offset) => {
    const response = myAxiosWithoutHeaders.get(`/chat/${userToChat}?offset=${offset}`)
    return response
}

export { sendProfileRequest, sendLoginRequest, sendLogoutRequest, sendUsersOnlineRequest, sendChatMessagesRequest, sendChatMessagesRequestWithOffset, sendGroupChatMessagesRequest };
