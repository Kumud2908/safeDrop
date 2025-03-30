import React, { useState } from "react";
import CryptoJS from "crypto-js";

export default function Download() {
  const [id, setId] = useState("");
  const [key, setKey] = useState("");
  const [decryptedText, setDecryptedText] = useState(null);
  const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";


  const handleDecrypt = async () => {
    try {
        if (!id || !key) {
            setDecryptedText("Please enter ID and Key.");
            return;
        }

        console.log("id: ", id);
        const response = await fetch(`${id}`);

        if (!response.ok) throw new Error("Failed to fetch data");

        const encryptedArrayBuffer = await response.arrayBuffer();
        console.log("encryptedText: ", encryptedArrayBuffer);
        const encryptedBytes = new Uint8Array(encryptedArrayBuffer);

        if (encryptedBytes.length < 16) throw new Error("Invalid encrypted data");

        // Extract IV (first 16 bytes)
        const ivHex = encryptedBytes.slice(0, 16);
        const encryptedData = encryptedBytes.slice(16);

        // Convert IV to Hex String
        const ivHexString = Array.from(ivHex)
          .map(b => b.toString(16).padStart(2, "0"))
          .join("");

        const ivWordArray = CryptoJS.enc.Hex.parse(ivHexString);

        // Convert key from hex string to WordArray
        const keyWordArray = CryptoJS.enc.Hex.parse(key.trim());

        // Convert encrypted data to WordArray
        const encryptedWordArray = CryptoJS.lib.WordArray.create(encryptedData);

        // Perform decryption
        const decrypted = CryptoJS.AES.decrypt(
            { ciphertext: encryptedWordArray },
            keyWordArray,
            { iv: ivWordArray, mode: CryptoJS.mode.CTR, padding: CryptoJS.pad.NoPadding }
        );

        // Debug output: Decrypted as raw bytes (Hex)
        console.log("decrypted raw bytes:", decrypted);
        console.log("decrypted Hex:", decrypted.toString(CryptoJS.enc.Hex));

        // Convert decrypted Hex to WordArray (raw bytes)
        const decryptedBytes = CryptoJS.enc.Hex.parse(decrypted.toString(CryptoJS.enc.Hex));
        const decryptedArrayBuffer = new Uint8Array(decryptedBytes.words.length * 4);

        // Convert decrypted WordArray to Uint8Array
        for (let i = 0; i < decryptedBytes.words.length; i++) {
            const word = decryptedBytes.words[i];
            decryptedArrayBuffer[i * 4] = (word >> 24) & 0xff;
            decryptedArrayBuffer[i * 4 + 1] = (word >> 16) & 0xff;
            decryptedArrayBuffer[i * 4 + 2] = (word >> 8) & 0xff;
            decryptedArrayBuffer[i * 4 + 3] = word & 0xff;
        }

        console.log("Decrypted Uint8Array:", decryptedArrayBuffer);

        // Try to decode the raw bytes to UTF-8 text
        const decoder = new TextDecoder("utf-8", { fatal: true });
        let decryptedText = decoder.decode(decryptedArrayBuffer);

        // Remove trailing null characters (if any)
        decryptedText = decryptedText.replace(/\0/g, '').trim(); // Remove padding or unwanted null characters

        console.log("Decrypted text:", decryptedText);
        setDecryptedText(decryptedText || "Decryption failed");
    } catch (error) {
        console.error("Decryption Error:", error);
        setDecryptedText("Error decrypting data");
    }
};




  return (
    <div>
      <h2>Decrypt Data</h2>
      <input type="text" placeholder="Enter ID" value={id} onChange={(e) => setId(e.target.value)} />
      <input type="text" placeholder="Enter Key" value={key} onChange={(e) => setKey(e.target.value)} />
      <button onClick={handleDecrypt}>Decrypt</button>
      <h3>Decrypted Text:</h3>
  {decryptedText && (
    // If decryptedText is a string, display it as a single paragraph
    typeof decryptedText === "string" ? (
      <p>{decryptedText}</p>
    ) : (
      // If decryptedText is an array, map over it
      decryptedText.map((element, index) => <p key={index}>{element}</p>)
    )
  )}
    </div>
  );
}
