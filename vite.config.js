import { defineConfig } from 'vite';
import path from 'path';

export default defineConfig({
  root: path.resolve(__dirname, 'src/templates'),
  build: {
    outDir: path.resolve(__dirname, 'src/templates/dist'),  // Output directory inside 'templates/js'
    rollupOptions: {
      input: path.resolve(__dirname, 'src/templates/index.html'),
    },
    emptyOutDir: true,
  },
});
