import React, { useContext } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { UserContext } from './context/UserContext';
import './assets/styles/Home.css';

const BookLogo = () => (
  <svg width="32" height="32" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg" style={{marginRight: 12}}>
    <rect x="4" y="6" width="24" height="20" rx="4" fill="#fff" stroke="#222" strokeWidth="2"/>
    <rect x="8" y="10" width="16" height="12" rx="2" fill="#ededed" stroke="#222" strokeWidth="1.5"/>
    <line x1="16" y1="10" x2="16" y2="22" stroke="#222" strokeWidth="1.5"/>
  </svg>
);

function Toolbar() {
  const { user } = useContext(UserContext);
  const navigate = useNavigate();
  const isAdmin = user && user.user_type && user.user_type.toLowerCase() === 'administrador';

  const handleLogout = () => {
    localStorage.removeItem('userId');
    localStorage.removeItem('usertype');
    localStorage.removeItem('token');
    navigate('/login');
  };

  return (
    <nav className="toolbar">
      <div style={{display: 'flex', alignItems: 'center'}}>
        <BookLogo />
        <span className="toolbar-title">Portal de Cursos</span>
      </div>
      <div className="toolbar-nav">
        <Link to="/home" className="toolbar-link">Home</Link>
        <Link to="/my-courses" className="toolbar-link">Mis Cursos</Link>
        <Link to="/search" className="toolbar-link">Buscar</Link>
        {isAdmin && <Link to="/manage-courses" className="toolbar-link">Gestión de Cursos</Link>}
        {isAdmin && <Link to="/services-status" className="toolbar-link">Estado de Servicios</Link>}
      </div>
      <div style={{display: 'flex', alignItems: 'center'}}>
        {user && <span className="toolbar-user">{user.username}</span>}
        <button className="toolbar-logout" onClick={handleLogout}>Cerrar sesión</button>
      </div>
    </nav>
  );
}

export default Toolbar; 