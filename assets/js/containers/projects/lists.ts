import { connect } from 'react-redux'

import mapStateToProps from '../mapState'
import lists from '../../components/projects/lists.tsx'

export default connect(mapStateToProps)(lists)
