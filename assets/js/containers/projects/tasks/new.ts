import { connect } from 'react-redux'
import mapStateToProps from '../../mapState'
import task from '@/components/projects/tasks/new.tsx'

export default connect(mapStateToProps)(task)
