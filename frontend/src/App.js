import React from 'react';
import { BrowserRouter as Router, Route, Routes, Outlet } from 'react-router-dom';
import Home from './components/courses/Home';
import Login from './components/auth/Login';
import Register from './components/auth/Register';
import CourseDetails from './components/courses/CourseDetails';
import MyCourses from './components/courses/MyCourses';
import UploadFile from './components/courses/UploadFile';
import CourseFiles from './components/courses/CourseFiles';
import SearchCourses from './components/courses/SearchCourses';
import ProtectedRoute from './components/ProtectedRoute';
import AdminRoute from './components/AdminRoute';
import ManageCourses from './components/courses/ManageCourses';
import AddCourse from './components/courses/AddCourse';
import EditCourse from './components/courses/EditCourse';
import CourseComments from './components/courses/CourseComments'; 
import CommentForm from './components/courses/CommentForm';
import ServiceStatus from './components/admin/ServiceStatus';
import { Navigate } from 'react-router-dom';
import Toolbar from './components/Toolbar';
import { UserProvider } from './components/context/UserContext';

function WithToolbar() {
  return (
    <>
      <Toolbar />
      <Outlet />
    </>
  );
}

function App() {
  return (
    <UserProvider>
      <Router>
        <Routes>
          <Route path="/" element={<Navigate replace to="/login" />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          
          <Route element={<ProtectedRoute />}>
            <Route element={<WithToolbar />}>
              <Route path="/home" element={<Home />} />
              <Route path="/courses/:courseId" element={<CourseDetails />} />
              <Route path="/my-courses" element={<MyCourses />} />
              <Route path="/search" element={<SearchCourses />} />
              <Route path="/upload/:courseId" element={<UploadFile />} />
              <Route path="/courses/:courseId/comments" element={<CourseComments />} />
              <Route path="/courses/:courseId/files" element={<CourseFiles />} />
              <Route path="/courses/:courseId/comment" element={<CommentForm />} />
            </Route>
          </Route>

          <Route element={<AdminRoute />}>
            <Route element={<WithToolbar />}>
              <Route path="/manage-courses" element={<ManageCourses />} />
              <Route path="/add-course" element={<AddCourse />} />
              <Route path="/edit-course/:courseId" element={<EditCourse />} />
              <Route path="/services-status" element={<ServiceStatus />} />
            </Route>
          </Route>
        </Routes>
      </Router>
    </UserProvider>
  );
}

export default App;
