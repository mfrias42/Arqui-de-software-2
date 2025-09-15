import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import '../assets/styles/SearchCourses.css';

const SearchCourses = () => {
    const [searchTerm, setSearchTerm] = useState('');
    const [courses, setCourses] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const navigate = useNavigate();

    const handleSearch = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        try {
            const response = await axios.get(`http://localhost:8082/search?q=${searchTerm}`);
            if (!response.data || response.data.length === 0) {
                setError('No se encontraron cursos con ese nombre.');
                setCourses([]);
            } else {
                setCourses(response.data);
            }
        } catch (err) {
            setError('Error fetching courses: ' + err.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="search-courses-outer">
            <div className="search-courses-card">
            <button className="back-button" onClick={() => navigate('/home')}>Volver</button>
            <h1>Buscar Cursos</h1>
            <form onSubmit={handleSearch} className="search-form">
                <input
                    type="text"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    placeholder="Ingrese el nombre del curso"
                    required
                    className="search-input"
                />
                <button type="submit" className="search-button">Buscar</button>
            </form>
            {loading && <p>Cargando...</p>}
            {error && <p className="error-message">{error}</p>}
                {courses.length > 0 && (
                    <div className="search-courses-grid">
                    {courses.map(course => (
                            <div key={course.id} className="search-course-card">
                            <h3>{course.name}</h3>
                            <p>{course.description}</p>
                            <div className="my-course-actions">
                                <button className="details-button" onClick={() => navigate(`/courses/${course.id}`)}>Ver detalles</button>
                                <button className="upload-button" onClick={() => navigate(`/upload/${course.id}`)}>Subir Archivo</button>
                                <button className="comment-button" onClick={() => navigate(`/courses/${course.id}/comments`)}>Comentar</button>
                            </div>
                            </div>
                    ))}
                    </div>
            )}
                {(!loading && !error && courses.length === 0) && <p>No se encontraron cursos.</p>}
            </div>
        </div>
    );
};

export default SearchCourses;
