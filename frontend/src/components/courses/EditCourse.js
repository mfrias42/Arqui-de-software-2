import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import { useParams, useNavigate } from 'react-router-dom';
import '../assets/styles/EditCourse.css';

function EditCourse() {
    const { courseId } = useParams();
    const navigate = useNavigate();
    const [courseData, setCourseData] = useState({
        id: parseInt(courseId),
        name: '',
        description: '',
        category: '',
        duration: '',
        instructor_id: '',
        capacity: '',
        imageBase64: '',
    });
    const [error, setError] = useState('');
    const [imagePreview, setImagePreview] = useState('');
    const [imageName, setImageName] = useState('');
    const fileInputRef = useRef();

    useEffect(() => {
        const fetchCourseData = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/courses/${courseId}`);
                setCourseData({
                    id: response.data.id,
                    name: response.data.name,
                    description: response.data.description,
                    category: response.data.category,
                    duration: response.data.duration,
                    instructor_id: response.data.instructor_id,
                    capacity: response.data.capacity,
                    imageBase64: response.data.imageBase64 || '',
                });
                setImagePreview(response.data.imageBase64 || '');
                setImageName('');
            } catch (err) {
                setError('Error fetching course details: ' + err.message);
            }
        };
        fetchCourseData();
    }, [courseId]);

    const handleChange = (e) => {
        setCourseData({ ...courseData, [e.target.name]: e.target.value });
    };

    const handleImageChange = (e) => {
        const file = e.target.files[0];
        if (file) {
            setImageName(file.name);
            const reader = new FileReader();
            reader.onloadend = () => {
                setCourseData(prev => ({ ...prev, imageBase64: reader.result }));
                setImagePreview(reader.result);
            };
            reader.readAsDataURL(file);
        }
    };

    const handleFileButtonClick = (e) => {
        e.preventDefault();
        fileInputRef.current.click();
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        const token = localStorage.getItem('token');
        try {
            const payload = {
                ...courseData,
                instructor_id: Number(courseData.instructor_id),
                capacity: Number(courseData.capacity)
            };
            await axios.put(`http://localhost:8080/courses/${courseId}`, payload, {
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                }
            });
            alert('Curso actualizado con éxito');
            navigate('/manage-courses');
        } catch (error) {
            setError(`Error updating course: ${error.response?.data?.message || error.message}`);
        }
    };
    
    return (
        <div className="edit-course-bg">
            <div className="edit-course-container">
                <button className="back-button" onClick={() => navigate('/home')}>Volver</button>
                <h1>Editar Curso</h1>
                {error && <p className="error-message">{error}</p>}
                <form onSubmit={handleSubmit} className="edit-course-form">
                    <label>Nombre del curso</label>
                    <input name="name" value={courseData.name} onChange={handleChange} />

                    <label>Descripción</label>
                    <input name="description" value={courseData.description} onChange={handleChange} />

                    <label>Categoría</label>
                    <input name="category" value={courseData.category} onChange={handleChange} />

                    <label>Duración</label>
                    <input name="duration" value={courseData.duration} onChange={handleChange} />

                    <label>ID del Instructor</label>
                    <input 
                        name="instructor_id" 
                        type="number" 
                        value={courseData.instructor_id} 
                        onChange={handleChange} 
                        min="0" 
                    />

                    <label>Capacidad</label>
                    <input 
                        name="capacity" 
                        type="number" 
                        value={courseData.capacity} 
                        onChange={handleChange} 
                        min="1"
                    />

                    <label>Imagen del curso (opcional):</label>
                    <input 
                        type="file" 
                        accept="image/*" 
                        onChange={handleImageChange} 
                        ref={fileInputRef}
                        style={{display: 'none'}}
                    />
                    <button 
                        type="button"
                        className="file-upload-btn"
                        onClick={handleFileButtonClick}
                        style={{marginBottom: 8, display: 'block'}}
                    >
                        Seleccionar imagen
                    </button>
                    {imageName && <div style={{fontSize: '0.95rem', color: '#444', marginBottom: 8}}>{imageName}</div>}
                    {imagePreview && (
                        <img src={imagePreview} alt="Preview" style={{maxWidth: 180, margin: '10px auto', borderRadius: 8}} />
                    )}

                    <button type="submit" className="submit-button">Actualizar Curso</button>
                </form>
            </div>
        </div>
    );
}

export default EditCourse;
