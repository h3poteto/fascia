import * as listActions from '../actions/ListAction';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import ListView from '../components/ListView.jsx';
import mapStateToProps from './mapStateToProps';

function mapDispatchToProps(dispatch) {
  return bindActionCreators(listActions, dispatch);
}


export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ListView);
