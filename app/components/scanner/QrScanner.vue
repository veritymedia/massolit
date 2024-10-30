<template>
  <div class="flex items-center justify-center">
    <div
      id="qr-scanner-container"
      :style="scannerContainerSize"
      class="w-full aspect-square"
    >
      <div id="qr-scanner-box" class="w-full h-full border" ref="qrScannerBox">
        <div id="scanner-overlay"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref } from "vue";
import { Html5Qrcode } from "html5-qrcode";
import { ClientResponseError } from "pocketbase";

const emit = defineEmits(["qrResult"]);

const { width: windowWidth, height } = useWindowSize();
const scannerContainerSize = computed(() => {
  return { width: windowWidth.value - 100 };
});

const qrScannerBox = ref(null);
let html5QrCode: Html5Qrcode | null = null;

const scanResultText = ref("");

const createScanner = () => {
  html5QrCode = new Html5Qrcode("qr-scanner-box");
};

function startScanner() {
  if (!html5QrCode) {
    return;
  }
  html5QrCode
    .start(
      { facingMode: "environment" }, // Use back camera
      {
        fps: 20, // Frames per second
        aspectRatio: 1,
      },
      (decodedText, decodedResult) => {
        console.log(`QR Code detected: ${decodedText}`);
        scanResultText.value = decodedText;
        emit("qrResult", decodedText);
      },
      (errorMessage) => {
        console.log(`QR Code scan error: ${errorMessage}`);
      },
    )
    .catch((err) => {
      console.error(`Unable to start the QR Code scanner: ${err}`);
    });
}

const stopScanner = () => {
  if (html5QrCode) {
    html5QrCode
      .stop()
      .then(() => {
        console.log("QR Code scanning stopped.");
      })
      .catch((err) => {
        console.error(`Unable to stop scanning: ${err}`);
      });
  }
};

// Start scanner when component is mounted
onMounted(() => {
  nextTick(() => {
    createScanner();
    startScanner();
  });
});

onBeforeUnmount(() => {
  stopScanner();
});
</script>

<style scoped>
#scanner-overlay {
  position: absolute;
  top: 0;
  left: 0;
  border: 2px;
  width: 100%;
  height: 100%;
  background: rgba(5, 5, 5, 0.7); /* Overlay effect */
  pointer-events: none;
}
</style>
