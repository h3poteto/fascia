export default function mapStateToProps(state) {
  const { ProjectReducer } = state;
  return {
    routerState: state.router,
    ProjectReducer
  };
}
