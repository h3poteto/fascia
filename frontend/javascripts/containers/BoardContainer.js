import * as boardActions from '../actions/BoardAction';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import BoardView from '../components/BoardView.jsx';

function mapStateToProps(state) {
  const { BoardReducer } = state;
  return BoardReducer;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(boardActions, dispatch);
}


export default connect(
  mapStateToProps,
  mapDispatchToProps
)(BoardView);
