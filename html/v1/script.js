class HexGridSystem {
  constructor() {
    this.maxLevel = 6; // Up to 2^8 = 256x256
    this.currentLevel = 2;
    this.init();
  }

  init() {
    this.addLevelControls();
    this.generateGrid(this.currentLevel);
  }

  addLevelControls() {
    const container = document.querySelector('.container');
    const controls = document.createElement('div');
    controls.className = 'controls';
    controls.innerHTML = `
            <div style="text-align: center; margin: 20px 0;">
                <label for="level-select">Grid Level (2^(x+2) boxes): </label>
                <select id="level-select">
                    <option value="1">Level 1 (2x2 = 4)</option>
                    <option value="2" selected>Level 2 (4x4 = 16)</option>
                    <option value="3">Level 3 (8x8 = 64)</option>
                    <option value="4">Level 4 (16x16 = 256)</option>
                    <option value="5">Level 5 (32x32 = 1024)</option>
                    <option value="6">Level 6 (64x64 = 4096)</option>
                </select>
                <button id="generate-btn">Generate Grid</button>
            </div>
        `;
    container.insertBefore(controls, container.querySelector('h1').nextSibling);

    document.getElementById('generate-btn').addEventListener('click', () => {
      const level = parseInt(document.getElementById('level-select').value);
      this.generateGrid(level);
    });
  }

  generateGrid(level) {
    // Clear existing grids
    const existingGrids = document.querySelectorAll('.grid-level');
    existingGrids.forEach(grid => grid.remove());

    const gridSize = Math.pow(2, level);
    const totalBoxes = gridSize * gridSize;

    const gridContainer = document.createElement('div');
    gridContainer.className = 'grid-level';
    gridContainer.id = `level-${level}`;

    // Add CSS for this level
    this.addLevelCSS(level, gridSize);

    // Create array to hold boxes in Z-order
    const boxes = new Array(totalBoxes);

    // Generate boxes with proper Z-order hex values
    for (let i = 0; i < totalBoxes; i++) {
      const hexBox = document.createElement('div');
      hexBox.className = 'hex-box';

      // Convert Z-order index to grid position
      const { row, col } = this.zOrderToPosition(i, gridSize);
      const gridIndex = row * gridSize + col;

      const hexValue = i.toString(16).padStart(4, '0').toLowerCase();
      hexBox.textContent = hexValue;

      hexBox.addEventListener('click', () => {
        console.log(`Clicked box: ${hexValue} (Z-order: ${i}, Grid pos: ${row},${col})`);
        hexBox.classList.toggle('completed');
      });

      boxes[gridIndex] = hexBox;
    }

    // Add boxes to container in grid order
    boxes.forEach(box => gridContainer.appendChild(box));
    document.querySelector('.container').appendChild(gridContainer);
  }

  zOrderToPosition(zIndex, gridSize) {
    let x = 0, y = 0;
    let level = Math.log2(gridSize);

    // Process each bit level from most significant to least
    for (let i = level - 1; i >= 0; i--) {
      const mask = 1 << i;
      const bit0 = (zIndex >> (2 * i)) & 1;     // x bit
      const bit1 = (zIndex >> (2 * i + 1)) & 1; // y bit

      x |= bit0 << i;
      y |= bit1 << i;
    }

    return { row: y, col: x };
  }

  addLevelCSS(level, gridSize) {
    const styleId = `level-${level}-style`;
    let existingStyle = document.getElementById(styleId);
    if (existingStyle) {
      existingStyle.remove();
    }

    const style = document.createElement('style');
    style.id = styleId;

    const boxSize = Math.max(20, 500 / gridSize);
    const fontSize = Math.max(6, boxSize / 4);

    style.textContent = `
            #level-${level} {
                display: grid;
                grid-template-columns: repeat(${gridSize}, 1fr);
                gap: 1px;
                max-width: ${gridSize * (boxSize + 1)}px;
                margin: 0 auto;
            }
            
            #level-${level} .hex-box {
                min-width: ${boxSize}px;
                min-height: ${boxSize}px;
                font-size: ${fontSize}px;
            }
        `;

    document.head.appendChild(style);
  }
}

// Initialize the grid system when the page loads
document.addEventListener('DOMContentLoaded', () => {
  new HexGridSystem();
});