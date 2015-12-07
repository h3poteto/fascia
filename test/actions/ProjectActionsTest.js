import * as projectActions from '../../frontend/javascripts/actions/ProjectAction';
import expect from 'expect';

describe('closeFlash', () => {
  it('should close flash', () => {
    const expectedAction = {
      type: projectActions.CLOSE_FLASH
    };
    expect(projectActions.closeFlash()).toEqual(expectedAction);
  });
});
