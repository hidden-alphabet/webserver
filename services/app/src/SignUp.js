import { Link } from 'react-router-dom'
import React from 'react'

import { Fontified } from './Dashboard'
import { Header } from './Header'

const Field = props =>
  <div>
  <p style={{ margin: '0' }}>{props.name}:</p>
  <input></input>
  </div>

const CreateAccountButton = _ =>
  <div style={{ color: "white", paddingRight: "20px" }}>
    <Link to="/dashboard" style={{ display: 'inline-block' }}>
      <input
      id='call-to-action'
      type="button" 
      style={{ 
        color: "white", 
        width: '200px', 
        height: '40px',
        border: 'none',
        borderRadius: '10px 10px 10px 10px', 
        backgroundColor: '#1aab81',
        boxShadow: '2px',
        fontSize: '13px'
      }} 
      value="Create Account"
      >
      </input>
    </Link>
  </div>

export class Signup extends React.Component {
  render() {
    return(
      <Fontified>
          <Header />
          <div style={{ borderRadius: '4px', boxShaddow: '0 15px 35px 0 rgba(60,66,87, 0.1), 0 5px 15px 0 rgba(0, 0, 0, .07)', height: '100%' }}>
            <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'center', justifyContent: 'space-around' }}>
              <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'space-around' }}>
                <Field name="Full Name" />
                <Field name="Email" />
                <Field name="Password" />
                <Field name="Confirm Password" />
                <div style={{ padding: '10px' }}>
                  <CreateAccountButton />
                </div>
              </div>
            </div>
          </div>
      </Fontified>
    )
  }
}