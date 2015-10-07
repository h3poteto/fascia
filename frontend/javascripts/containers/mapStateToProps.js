export default function mapStateToProps(state) {
  const { BoardReducer } = state;
  return {
    routerState: state.router,
    BoardReducer
  };
}
