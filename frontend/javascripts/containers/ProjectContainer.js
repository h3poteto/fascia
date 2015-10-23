import * as projectActions from '../actions/ProjectAction';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import ProjectView from '../components/ProjectView.jsx';
import mapStateToProps from './mapStateToProps';

function mapDispatchToProps(dispatch) {
  return bindActionCreators(projectActions, dispatch);
}


export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ProjectView);
