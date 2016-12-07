import * as menuActions from '../actions/MenuAction'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import MenuView from '../components/MenuView.jsx'
import mapStateToProps from './mapStateToProps'

function mapDispatchToProps(dispatch) {
  return bindActionCreators(menuActions, dispatch)
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(MenuView)
