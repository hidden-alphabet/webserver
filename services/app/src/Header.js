import { Link } from 'react-router-dom' 
import React from 'react'

export const Header = props =>
  <header style={{ 
    backgroundColor: "black", 
    height: "70px", 
    opacity: '0%',
    color: 'white',
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingLeft: '172px',
    paddingRight: '172px'
  }}>
    <Link to='/'>
        <p 
        style={{ 
            paddingLeft: '30px', 
            fontSize: '20',
            color: 'white'
        }}>
            Hidden Alphabet
        </p>
    </Link>
    {props.children}
  </header>