import axios from 'axios';

const registrationData = {
  email: "parvaggarwal130@gmail.com",
  name: "Parv Aggarwal",
  mobileNo: "9050547124",
  githubUsername: "parvaggarwal01",
  rollNo: "12310700",
  accessCode: "TRvZWq",
  clientID: '75b8b764-2c61-4388-aa96-d3c126ec1805',
  clientSecret: 'HwxRATehYrvjbUDN'
};

async function login() {
  try {
    const res = await axios.post(
      "http://4.224.186.213/evaluation-service/auth",
      registrationData,
    );
    console.log("Registration Successful! SAVE THESE CREDENTIALS:");
    console.log(res.data);
  } catch (err) {
    console.error("Failed:", err.response?.data || err.message);
  }
}

login();
