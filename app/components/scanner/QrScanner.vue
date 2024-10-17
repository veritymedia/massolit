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

<script setup>
import { onMounted, onBeforeUnmount, ref } from "vue";
import { Html5Qrcode } from "html5-qrcode";

const { width, height } = useWindowSize();
const scannerContainerSize = computed(() => {
  return { width: width - 100 };
});

const qrScannerBox = ref(null);
let html5QrCode = null;
const scanResultText = ref("");
const startScanner = () => {
  html5QrCode = new Html5Qrcode("qr-scanner-box");

  html5QrCode
    .start(
      { facingMode: "environment" }, // Use back camera
      {
        fps: 20, // Frames per second
        aspectRatio: 1,

        // qrbox: { width: 200, height: 200 }, // Define scanner box size
      },
      (decodedText, decodedResult) => {
        console.log(`QR Code detected: ${decodedText}`);
        scanResultText.value = decodedText;
      },
      (errorMessage) => {
        console.log(`QR Code scan error: ${errorMessage}`);
      },
    )
    .catch((err) => {
      console.error(`Unable to start the QR Code scanner: ${err}`);
    });
};

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
    startScanner();
  });
});

// Stop scanner when component is destroyed
onBeforeUnmount(() => {
  stopScanner();
});

function checkCodeValidity() {
  // check if code exists in database
  // if yes, check if it is withdrawn by someone
  //    if yes return user, navigateTo "rental" page which shows book and user who took it
  //    if no, navigateTo page to search for managebac user.
  // if no:
  //    ask if should create book in database.
  //        if yes, add book to db, go back to scan screen.
  //        if no, go back to scan screen.
}
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
