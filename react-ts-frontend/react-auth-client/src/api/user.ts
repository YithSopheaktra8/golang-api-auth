import api from "./axios";


export async function getProfile(){

    const response =
        await api.get(
            "/user/me"
        );


    return response.data;

}