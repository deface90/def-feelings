const logger = store => next => action => {
  //console.group(action.type)
  //console.info('dispatching', action)
  //console.groupEnd()
  return next(action)
}

export default logger
