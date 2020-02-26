import { connect } from 'react-redux'

import mapStateToProps from './mapState'
import projects from '../components/projects.tsx'

export default connect(mapStateToProps)(projects)
