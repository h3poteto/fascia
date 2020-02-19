import * as settingsActions from '../actions/UserSettings'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import UserSettings from '../components/UserSettings.jsx'
import mapStateToProps from './mapStateToProps'

function mapDispatchToProps(dispatch) {
  return {
    settingsActions: bindActionCreators(settingsActions, dispatch)
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(UserSettings)
