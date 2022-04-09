const setLocalStorage = function (key, value) {
    localStorage.setItem(key, value);
}
const getLocalStorage = function (key) {
    return localStorage.getItem(key)
}
const clearLocalStorage = function (key) {
    localStorage.removeItem(key);
}

export { setLocalStorage, getLocalStorage, clearLocalStorage }
