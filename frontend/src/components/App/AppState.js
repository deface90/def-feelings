const UPDATE_USER = 'app/UPDATE_USER'
const UPDATE_NOTIFICATION = 'app/UPDATE_NOTIFICATION'

const initialState = {
    user: {},
    notification: [],
};

export function setUser(user) {
    return {
        type: UPDATE_USER,
        user
    }
}

export function setNotification(notification) {
    return {
        type: UPDATE_NOTIFICATION,
        notification
    }
}

// Reducer
export default function AppState(state = initialState, action = {}) {
    switch (action.type) {
        case UPDATE_USER: {
            let user = action.user

            return {
                ...state,
                user
            }
        }
        case UPDATE_NOTIFICATION: {
            let notification = action.notification

            return {
                ...state,
                notification
            }
        }
        default:
            return state;
    }
}
