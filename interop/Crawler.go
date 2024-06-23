package interop

/*
#cgo CFLAGS: -I/usr/include/python3.11 -I/usr/lib/jvm/java-11-openjdk-amd64/include -I/usr/lib/jvm/java-11-openjdk-amd64/include/linux
#cgo LDFLAGS: -lpython3.11

#include <Python.h>
#include <stdlib.h>
#include <stdio.h>

// The rest of the C code

char* load_cpython() {
	if(Py_IsInitialized())
	{
	return NULL;}

	Py_Initialize();
	if (!Py_IsInitialized()) {
		return "Failed to initialize Python. Py_Initialize did not set up Python interpreter.";}

	// Get the “path” attribute from “sys” module in Python
	PyObject* sysPath = PySys_GetObject("path");
	if (sysPath == NULL) {
		return "Failed to get sys.path";
	}

	// Create a Python string representing the current directory
	PyObject* cwd = PyUnicode_DecodeFSDefault(".");
	if (cwd == NULL) {
		return "Failed to decode current directory";
	}
	// Append the current directory to the sys.path list
	if (PyList_Append(sysPath, cwd) != 0) {
		return "Failed to append current directory to sys.path";
	}
	Py_DECREF(cwd);
	return NULL;
}
char* handle_error()
{
	// Fetch the error type, value, and traceback
	PyObject *type, *value, *traceback;
	// fills the pointers with the error information
	PyErr_Fetch(&type, &value, &traceback);
	// Get the name of the error type
	const char* error_name = PyExceptionClass_Name(type);

	// Get the string representation of the error value
	PyObject* value_of_error_obj = PyObject_Str(value);
	PyObject* bytes_utf8_value = PyUnicode_AsUTF8String(value_of_error_obj);
	char* value_as_c_string = PyBytes_AsString(bytes_utf8_value);

	// Get the string representation of the traceback
	PyObject* traceback_str = PyObject_Str(traceback);
	PyObject* bytes_utf8_traceback = PyUnicode_AsUTF8String(traceback_str);
	char* traceback_as_c_string = PyBytes_AsString(bytes_utf8_traceback);

	// allocate memory for the error message (4 = 1 colon + 1 space + 1 newline + 1 null terminator)
	char* res = malloc(strlen(error_name) + strlen(value_as_c_string) + strlen(traceback_as_c_string) + 4);
	// format the message
	sprintf(res, "%s: %s\n%s", error_name, value_as_c_string, traceback_as_c_string);

	// Release the memory allocated by Python
	Py_DECREF(type);
	Py_DECREF(value);
	Py_DECREF(traceback);
	Py_DECREF(value_of_error_obj);
	Py_DECREF(traceback_str);
	// Clear the error Python has raised
	PyErr_Clear();
	return res;
}
// Get links from URL
char** extract_links_from_url(char* url, int depth, char** out_error)
{
	// initialize the variables
	PyObject* pName = NULL;
	PyObject* pModule = NULL;
	PyObject* pFunc = NULL;
	PyObject* pArgs = NULL;
	PyObject* pValue = NULL;
	char** result = NULL;
	// Ensure thread holds the GIL (Global Interpreter Lock)
	// Either the thread gets the GIL or it waits until it is available
	// Only a thread holding the GIL can use the interpreter
	// and execute Python code
	PyGILState_STATE gstate;
	gstate = PyGILState_Ensure();

	// create a python string object from the crawler module
	pName = PyUnicode_DecodeFSDefault("crawler");

	// import the module
	pModule = PyImport_Import(pName);

	// free the python string object
	Py_DECREF(pName);
	if(pModule == NULL) // handle error
	{
		// if Python error has been raised, lets get the error from there, otherwise, return a default error message
		*out_error = PyErr_Occurred() ?
				handle_error() :
					strdup("Failed to load module crawler");
		PyErr_Clear();
		goto cleanup;
	}

	// load the function extract_links_from_url from the module
	// functions are Python attributes of the module
	// function in Python is a callable object
	pFunc = PyObject_GetAttrString(pModule, "extract_links_from_url");
	if(!pFunc || !PyCallable_Check(pFunc)) // handle error
	{
		*out_error = PyErr_Occurred() ?
				handle_error() :
					strdup("Cannot find function extract_links_from_url");
		PyErr_Clear();
		goto cleanup;
	}

	pArgs = PyTuple_New(2); // create the tuple that holds the arguments
	PyTuple_SetItem(pArgs, 0, PyUnicode_FromString(url)); // url parameter
	PyTuple_SetItem(pArgs, 1, PyLong_FromLong(depth)); // depth parameter

	// call the function with the arguments and get the result
	pValue = PyObject_CallObject(pFunc, pArgs);

	// free the arguments tuple
	// tuple makes sure to decrement the reference count of the objects it holds
	Py_DECREF(pArgs);

	// check if error occurred during function execution
	if(PyErr_Occurred() || pValue == NULL)
	{
		*out_error = PyErr_Occurred() ?
			handle_error() :
				strdup("function extract_links_from_url failed");
		PyErr_Clear();
		goto cleanup;
	}

	// check if the result is a list
	if(!PyList_Check(pValue))
	{
		*out_error = strdup("function extract_links_from_url did not return a list");
		goto cleanup;
	}

	// copy the list to a C array of strings
	Py_ssize_t size = PyList_Size(pValue); // get the size of the list
	result = malloc((size + 1) * sizeof(char*));
	result[size] = NULL; // mark the last element using NULL


	// copy the strings from the list to the C array
	for (Py_ssize_t i = 0; i < size; i++)
	{
		PyObject *item = PyList_GetItem(pValue, i); // the i-th string
		// make sure it is a string
		if(!PyUnicode_Check(item))
		{
			*out_error = strdup("function extract_links_from_url returned a non-string item");
			// free the memory of already allocated strings
			for (Py_ssize_t j = 0; j < i; j++)
			{
				free(result[j]);
			}
			free(result);
			goto cleanup;
		}
		PyObject* item_as_utf8 = PyUnicode_AsUTF8String(item); // convert to bytes as utf-8
		result[i] = strdup(PyBytes_AsString(item_as_utf8)); // copy the bytes to a new string
		Py_DECREF(item_as_utf8); // free the bytes object
	}

	cleanup:
		Py_XDECREF(pFunc);
		Py_XDECREF(pModule);
		Py_XDECREF(pName);
		Py_XDECREF(pArgs);
		Py_XDECREF(pValue);
		PyGILState_Release(gstate);
		return result;
}



char* get_element(char** array, int index) {
	return array[index];
}

*/
import "C"
import (
	"errors"
	"unsafe"
)

func LoadPython() error {
	err := C.load_cpython()
	if err != nil {
		return errors.New(C.GoString(err))
	}
	return nil
}

func ExtractLinksFromURL(url string, depth int) ([]string, error) {
	c_url := C.CString(url)
	defer C.free(unsafe.Pointer(c_url))

	var c_error *C.char
	c_result := C.extract_links_from_url(c_url, C.int(depth), &c_error)
	if c_error != nil {
		defer C.free(unsafe.Pointer(c_error))
		return nil, errors.New(C.GoString(c_error))
	}

	defer C.free(unsafe.Pointer(c_result))
	length := 0
	for C.get_element(c_result, C.int(length)) != nil {
		length++
	}

	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(c_result))[:length:length]
	goStrings := make([]string, length)
	for i, s := range tmpslice {
		goStrings[i] = C.GoString(s)
		C.free(unsafe.Pointer(s))
	}

	return goStrings, nil
}
