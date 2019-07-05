import { Link } from 'react-router-dom'
import { Fontified } from './Dashboard'
import { Header } from './Header'
import React from 'react'

export class Login extends React.Component {
    render() {
        return (
            <Fontified>
                <Header />
                <div>
                    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', alignContent: 'center' }}>
                        <div>
                            <p style={{ margin: '0' }}>Email:</p>
                            <input></input>
                        </div>
                        <div>
                            <p style={{ margin: '0' }}>Password:</p>
                            <input></input>
                        </div>
                        <Link to='/dashboard'><input type="button" value="Login"></input></Link>
                    </div>
                </div>
            </Fontified>
        )
    }
}