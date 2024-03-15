function handler(req, res) {
    res.status(200).json({message: 'Hello, world!', params: req.query});
}

export default handler
