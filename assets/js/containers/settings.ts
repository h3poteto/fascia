import { connect } from 'react-redux'

import mapStateToProps from './mapState'
import settings from '../components/settings.tsx'

export default connect(mapStateToProps)(settings)
