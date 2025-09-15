import React, { useState, useContext } from 'react';
import axios from 'axios';
import { UserContext } from '../context/UserContext';
import { useNavigate } from 'react-router-dom';
import './Register.css'; 

function Register() {
    const [username, setUsername] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [user_type, setUserType] = useState('alumno');
    const [error, setError] = useState('');
    const { setUser } = useContext(UserContext);
    const navigate = useNavigate();

    const handleRegister = async (e) => {
        e.preventDefault();
        try {
            const response = await axios.post('http://localhost:8083/users', {
                username,
                email,
                password,
                user_type
            });
            setUser(response.data);
            alert('Registro exitoso');
            navigate('/login');
        } catch (error) {
            setError('Failed to register: ' + error.message);
        }
    };

    return (
        <div className="register-outer">
            <div className="register-card">
                <h1 className="register-title">Crear cuenta</h1>
                <p className="register-subtitle">Crea tu cuenta para acceder al portal</p>
            <form onSubmit={handleRegister} className="register-form">
                    <input type="text" value={username} onChange={e => setUsername(e.target.value)} placeholder="Usuario" required className="input-field" />
                <input type="email" value={email} onChange={e => setEmail(e.target.value)} placeholder="Email" required className="input-field" />
                    <input type="password" value={password} onChange={e => setPassword(e.target.value)} placeholder="Contraseña" required className="input-field" />
                <select value={user_type} onChange={e => setUserType(e.target.value)} required className="input-field">
                    <option value="alumno">Alumno</option>
                    <option value="administrador">Administrador</option>
                </select>
                    <button type="submit" className="register-button">Registrarse</button>
            </form>
                <div className="register-login-link">
                    <span>¿Ya tienes cuenta?</span>
                    <button type="button" className="login-link-btn" onClick={() => navigate('/login')}>Iniciar sesión</button>
                </div>
            {error && <p className="error-message">{error}</p>}
            </div>
        </div>
    );
}

export default Register;
