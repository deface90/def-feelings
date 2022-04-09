import {combineReducers} from 'redux';

import app from '../components/App/AppState';
import login from '../pages/Login/LoginState';
import register from '../pages/Register/RegisterState';

const rootReducer = combineReducers({
  app,
  login,
  register
});

export default rootReducer
