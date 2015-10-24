export default function mapStateToProps(state) {
  const { ProjectReducer, ListReducer } = state;
  return {
    routerState: state.router,
    ProjectReducer,
    ListReducer
  };
}
