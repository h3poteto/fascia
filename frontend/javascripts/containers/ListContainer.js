import * as listActions from '../actions/ListAction.js'
import * as newListModalActions from '../actions/ListAction/NewListModalAction.js'
import * as editListModalActions from '../actions/ListAction/EditListModalAction.js'
import * as newTaskModalActions from '../actions/ListAction/NewTaskModalAction.js'
import * as editProjectModalActions from '../actions/ListAction/EditProjectModalAction.js'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import ListView from '../components/ListView.jsx'
import mapStateToProps from './mapStateToProps'

function mapDispatchToProps(dispatch) {
  return bindActionCreators(Object.assign(
    {},
    listActions,
    newListModalActions,
    editListModalActions,
    newTaskModalActions,
    editProjectModalActions
  ), dispatch)
}


export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ListView)
