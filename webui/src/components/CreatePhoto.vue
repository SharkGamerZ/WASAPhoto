<script>
export default {
  name: 'CreatePhoto',
  props: {
    show: {
      type: Boolean,
      required: true
    }
  },
  data() {
    return {
      selectedFile: null,
      caption: '',
      loading: false,
      error: null,
      preview: null
    }
  },
  emits: ['close'],
  methods: {
    handleFileSelect(event) {
      const file = event.target.files[0];
      if (file) {
        this.selectedFile = file;
        // Create preview
        const reader = new FileReader();
        reader.onload = e => this.preview = e.target.result;
        reader.readAsDataURL(file);
      }
    },
    async uploadPhoto() {
      if (!this.selectedFile) return;
      
      this.loading = true;
      this.error = null;
      
      try {
        const formData = new FormData();
        formData.append('photo', this.selectedFile);
        formData.append('caption', this.caption);
        
        await this.$axios.post(`/users/${localStorage.getItem('token')}/photos`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        });
        
        this.$emit('close');
        // Optionally refresh the feed
        this.$router.go(0);
      } catch (e) {
        this.error = e.toString();
      }
      this.loading = false;
    },
    closeModal() {
      this.$emit('close');
      this.selectedFile = null;
      this.caption = '';
      this.preview = null;
      this.error = null;
    }
  }
}
</script>

<template>
  <transition name="modal">
    <div v-if="show" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>Create New Post</h2>
          <button class="close-button" @click="closeModal">
            <svg class="feather">
              <use href="/feather-sprite-v4.29.0.svg#x" />
            </svg>
          </button>
        </div>

        <div class="modal-body">
          <div v-if="error" class="error-message">
            <svg class="feather error-icon">
              <use href="/feather-sprite-v4.29.0.svg#alert-circle" />
            </svg>
            {{ error }}
          </div>
          
          <div class="upload-area" v-if="!preview">
            <input 
              type="file" 
              accept="image/*" 
              @change="handleFileSelect" 
              id="photo-upload" 
              class="hidden"
            />
            <label for="photo-upload" class="upload-label">
              <div class="upload-icon">
                <svg class="feather">
                  <use href="/feather-sprite-v4.29.0.svg#image" />
                </svg>
              </div>
              <span>Click to select a photo</span>
              <span class="upload-hint">JPG, PNG up to 10MB</span>
            </label>
          </div>

          <div v-else class="preview-area">
            <div class="preview-container">
              <img :src="preview" alt="Preview" class="photo-preview" />
              <button class="change-photo-button" @click="selectedFile = null; preview = null">
                Change photo
              </button>
            </div>
            <textarea 
              v-model="caption" 
              placeholder="Write a caption..." 
              class="caption-input"
              rows="3"
            ></textarea>
          </div>

          <button 
            @click="uploadPhoto" 
            class="upload-button" 
            :disabled="!selectedFile || loading"
          >
            <span v-if="!loading">Share post</span>
            <span v-else class="loading-text">
              <svg class="feather spinner">
                <use href="/feather-sprite-v4.29.0.svg#loader" />
              </svg>
              Uploading...
            </span>
          </button>
        </div>
      </div>
    </div>
  </transition>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  backdrop-filter: blur(5px);
}

.modal-content {
  background: white;
  border-radius: 16px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #eee;
}

.modal-header h2 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #1a1a1a;
}

.close-button {
  background: none;
  border: none;
  padding: 0.5rem;
  cursor: pointer;
  border-radius: 50%;
  transition: background-color 0.2s;
}

.close-button:hover {
  background-color: #f5f5f5;
}

.modal-body {
  padding: 1.5rem;
}

.upload-area {
  border: 2px dashed #e0e0e0;
  border-radius: 12px;
  padding: 2.5rem 1.5rem;
  text-align: center;
  transition: all 0.2s;
}

.upload-area:hover {
  border-color: #007bff;
  background: #f8f9ff;
}

.upload-label {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  cursor: pointer;
}

.upload-icon {
  width: 64px;
  height: 64px;
  background: #f0f4ff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 0.5rem;
}

.upload-icon .feather {
  width: 32px;
  height: 32px;
  stroke: #007bff;
}

.upload-hint {
  font-size: 0.875rem;
  color: #666;
}

.hidden {
  display: none;
}

.preview-area {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.preview-container {
  position: relative;
  border-radius: 12px;
  overflow: hidden;
  background: #f8f9fa;
}

.photo-preview {
  width: 100%;
  max-height: 400px;
  object-fit: contain;
  display: block;
}

.change-photo-button {
  position: absolute;
  bottom: 1rem;
  right: 1rem;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  cursor: pointer;
  font-size: 0.875rem;
  transition: background-color 0.2s;
}

.change-photo-button:hover {
  background: rgba(0, 0, 0, 0.8);
}

.caption-input {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  resize: vertical;
  font-size: 0.95rem;
  transition: border-color 0.2s;
}

.caption-input:focus {
  outline: none;
  border-color: #007bff;
}

.upload-button {
  width: 100%;
  padding: 0.875rem;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 500;
  transition: background-color 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  margin-top: 1rem;
}

.upload-button:hover:not(:disabled) {
  background: #0056b3;
}

.upload-button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.loading-text {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.spinner {
  animation: spin 1s linear infinite;
}

.error-message {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #dc3545;
  margin-bottom: 1rem;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  background: #fff5f5;
  font-size: 0.9rem;
}

.error-icon {
  stroke: #dc3545;
}

/* Modal transition animations */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style> 