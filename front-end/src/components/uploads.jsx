import React, {  useState } from "react";

export default function Uploads() {
  const [file, setFile] = useState(null);
//   const [preview, setPreview] = useState(null);
  const [uploading, setUploading] = useState(false);
 const [data,setData]=useState(null);
 const API_BASE_URL =  import.meta.env.VITE_API_URL

  const handleFileChange = (event) => {
    const selectedFile = event.target.files[0];
    if (selectedFile) {
      setFile(selectedFile);
    //   setPreview(URL.createObjectURL(selectedFile)); // ✅ Corrected createObjectURL
    }
  };

  const handleUpload = async () => {
    if (!file) {
      alert("Please select a file first!");
      return;
    }

    const formData = new FormData();
    formData.append("file", file);
    setUploading(true);

    try {
      const response = await fetch(`${API_BASE_URL}/upload`, {
        method: "POST",
        // ❌ Remove "Content-Type", let the browser set it
        body: formData,
       
      });

      if (response.ok) {
        const jsonData = await response.json(); // ✅ Parse response as JSON
        setData(jsonData)
        alert(`File uploaded successfully! Download link: ${jsonData.file_url}`);
      } else {
        alert("Upload failed. Please try again.");
      }
    } catch (error) {
      console.error("Upload error:", error);
      alert("Error uploading file");
    } finally {
      setUploading(false);
    }
  };

  return (
    <div style={{ padding: "20px", textAlign: "center" }}>
      <h2>Upload a File</h2>
      <input type="file" onChange={handleFileChange} />
      
      {/* {preview && (
        <div style={{ marginTop: "10px" }}>
          <p>File Preview:</p>
          <img src={preview} alt="Preview" style={{ maxWidth: "200px" }} />
        </div>
      )} */}

      <button onClick={handleUpload} disabled={uploading} style={{ marginTop: "10px" }}>
        {uploading ? "Uploading..." : "Upload File"}
      </button>
      {data && (
        <div style={{marginTop:"10px"}}>
            Data uploaded successfully
        </div>
      )}
      {
        data && (
            <div style={{marginTop:"10px"}}>
                 <p>URL:{data.file_url}</p>
                <p>Key:{data.key}</p>
               
            </div>
        )
      }
     
    </div>
  );
}
