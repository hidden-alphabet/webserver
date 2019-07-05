import { Link } from 'react-router-dom'
import React from 'react'

import { Header } from './Header'
import { Fontified } from './Dashboard'

export class Account extends React.Component {
  render() {
    return (
      <Fontified>
        <Header>
          <Link to="/dashboard">Back</Link>
        </Header>
        <hr/>
        <div style={{ display: "flex", flexDirection: "row", alignItems: "center" }}>
          <h2 style={{paddingRight: "5px"}}>Account</h2><a href='#'>edit</a>
        </div>
        <hr/>
        {"Full Name: Cole Hudson"}
        <br/>
        {"Email: cole+hiddenalphabet@colejhudson.com"}
        <br/>
        {"Password: ********"}
        <div style={{ display: "flex", flexDirection: "row", alignItems: "center" }}>
          <h2 style={{paddingRight: "5px"}}>API Keys</h2><a href='#'>Create Keys</a>
        </div>
        <hr/>
        {"Key: AKIAWYUWRXF4SADLT2LG"}{" "}<a href='#'>Copy</a>
        <br/>
        {"Secret: eROM6sMkzUMbmv3coJ00u66gtU2T6u0ARsSqKegp"}{" "}<a href='#'>Copy</a>
        <br/>
        <button>Download as CSV</button>
        <h3>Example</h3>
        <div>
          <hr/>
          <br/>
          <div style={{background: "lightgrey"}}>
          <code>
          {"1. curl api.hiddenalphabet.com \\"}<br/>
          {"2. -X POST \\"}<br/>
          {"3. --data '{ \"key\": \"AKIAWYUWRXF4SADLT2LG\", secret: \"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\" }\'"}<br/>
          </code>
          </div>
        </div>
      </Fontified>
    )
  }
}