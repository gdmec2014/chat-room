import { GetAllRoom } from "../../api/room";
import store from "../../store"

let room = {
    state: {
        data: {
            all_room: [],
            in_page: 0,
            open_room_data: {
                id: "",
                name: "",
                user_type: 0
            },
            ws_data: {
                event_type: 0,
                room: {
                    id: "",
                    name: "",
                    user_type: 0,
                    max_member: "2"
                },
                msg: "",
                data: {}
            }
        }
    },

    mutations: {
        SET_ALL_ROOM_DATA: (state, data) => {
            console.log("设置all room data ", data)
            state.data.all_room = data
        },
        SET_ROOM_PAGE: (state, in_page) => {
            console.log("设置all room page ", in_page)
            state.data.in_page = in_page
        },
        SET_OPEN_ROOM_DATA: (state, data) => {
            console.log("设置open room data ", data)
            state.data.open_room_data = data
        }
    },

    actions: {
        //获取所有的房间
        GetAllRoom({
            commit
        }) {
            GetAllRoom().then(res => {
                if (res.Result == 10000) {
                    commit('SET_ALL_ROOM_DATA', res.Data)
                } else {
                    //登陆失败处理
                    store.dispatch("SetWXUserData")
                }
            })
        },
        SetRoomPage({
            commit
        }, in_page) {
            commit('SET_ROOM_PAGE', in_page)
        },
        AddRoom({
            commit
        }, new_room) {
            let all_room = store.getters.room.all_room
            all_room.push(new_room)
            commit('SET_ALL_ROOM_DATA', all_room)
        },
        SetOpenRoomData({
            commit
        }, room) {
            let open_room_data = {
                id: room.id,
                name: room.name,
                user_type: 0
            }
            commit('SET_OPEN_ROOM_DATA', open_room_data)
        }
    }
}

export default room