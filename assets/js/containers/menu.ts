import { connect } from 'react-redux'
import { Dispatch, bindActionCreators } from 'redux'

import * as menuActions from '../actions/menu'
import menu from '../components/menu'
import mapStateToProps from './mapState'

const mapDispatchToProps = (dispatch: Dispatch) => {
  return {
    menuActions: bindActionCreators(menuActions, dispatch)
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(menu)
