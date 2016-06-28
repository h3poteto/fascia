import * as projectActions from '../actions/ProjectAction'
import * as newProjectModalActions from '../actions/ProjectAction/NewProjectModalAction'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import ProjectView from '../components/ProjectView.jsx'
import mapStateToProps from './mapStateToProps'

function mapDispatchToProps(dispatch) {
  return {
    projectActions: bindActionCreators(projectActions, dispatch),
    newProjectModalActions: bindActionCreators(newProjectModalActions, dispatch)
  }
}


export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ProjectView)
