import { connect } from 'react-redux'

import menu from '../components/menu.tsx'
import mapStateToProps from './mapState'

export default connect(mapStateToProps)(menu)
