import * as listActions from '../actions/ListAction.js'
import * as newListModalActions from '../actions/ListAction/NewListModalAction.js'
import * as editListModalActions from '../actions/ListAction/EditListModalAction.js'
import * as newTaskModalActions from '../actions/ListAction/NewTaskModalAction.js'
import * as editProjectModalActions from '../actions/ListAction/EditProjectModalAction.js'
import * as showTaskModalActions from '../actions/ListAction/ShowTaskModalAction.js'
import * as deleteProjectModalActions from '../actions/ListAction/DeleteProjectModalAction.js'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import ListView from '../components/ListView.jsx'
import mapStateToProps from './mapStateToProps'

function mapDispatchToProps(dispatch) {
  return {
    listActions: bindActionCreators(listActions, dispatch),
    newListModalActions: bindActionCreators(newListModalActions, dispatch),
    editListModalActions: bindActionCreators(editListModalActions, dispatch),
    newTaskModalActions: bindActionCreators(newTaskModalActions, dispatch),
    editProjectModalActions: bindActionCreators(editProjectModalActions, dispatch),
    showTaskModalActions: bindActionCreators(showTaskModalActions, dispatch),
    deleteProjectModalActions: bindActionCreators(deleteProjectModalActions, dispatch),
  }
}


export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ListView)
