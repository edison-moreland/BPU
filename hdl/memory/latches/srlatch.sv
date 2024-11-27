module srlatch(
    input  logic s, r,
    output logic q,
);
    always_ff @(posedge s, posedge r) begin
        if (r)
            q <= 0;
        else if(s)
            q <= 1;
    end
endmodule